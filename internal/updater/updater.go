package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/nadeemmhdm/sego-server/internal/config"
)

const repo = "nadeemmhdm/sego-server"

type Asset struct{ Name, BrowserDownloadURL string }
type Release struct {
	TagName string `json:"tag_name"`
	Draft bool `json:"draft"`
	Prerelease bool `json:"prerelease"`
	Assets []struct { Name string `json:"name"`; BrowserDownloadURL string `json:"browser_download_url"` } `json:"assets"`
}

type Info struct{ Current, Latest string; Available bool; Release Release }

func Latest(ctx context.Context, current string) (Info, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/"+repo+"/releases/latest", nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "SEGO-Server-Updater")
	resp, err := (&http.Client{Timeout: 15*time.Second}).Do(req)
	if err != nil { return Info{}, err }
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound { return Info{Current: normalize(current)}, nil }
	if resp.StatusCode != http.StatusOK { return Info{}, fmt.Errorf("GitHub release check returned HTTP %d", resp.StatusCode) }
	var rel Release
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&rel); err != nil { return Info{}, err }
	latest := normalize(rel.TagName); cur := normalize(current)
	return Info{Current: cur, Latest: latest, Available: newer(latest, cur), Release: rel}, nil
}

func normalize(v string) string { return strings.TrimPrefix(strings.TrimSpace(v), "v") }
func newer(a,b string) bool {
	pa, pb := strings.Split(normalize(a), "."), strings.Split(normalize(b), ".")
	for i:=0;i<3;i++ { var x,y int; if i<len(pa){x,_=strconv.Atoi(strings.FieldsFunc(pa[i],func(r rune)bool{return r<'0'||r>'9'})[0])}; if i<len(pb){y,_=strconv.Atoi(strings.FieldsFunc(pb[i],func(r rune)bool{return r<'0'||r>'9'})[0])}; if x!=y{return x>y} }
	return false
}

func ApplyLatest(ctx context.Context, current, configPath string) (bool, string, error) {
	info, err := Latest(ctx,current); if err != nil{return false,"",err}
	if !info.Available { return false, info.Latest, nil }
	assetName := packageName(info.Latest)
	var packageURL, sumsURL string
	for _,a := range info.Release.Assets { if a.Name==assetName {packageURL=a.BrowserDownloadURL}; if a.Name=="SHA256SUMS" {sumsURL=a.BrowserDownloadURL} }
	if packageURL=="" || sumsURL=="" { return false,info.Latest,fmt.Errorf("release %s is missing %s or SHA256SUMS",info.Release.TagName,assetName) }
	tmp, err := os.MkdirTemp("","sego-update-"); if err!=nil{return false,info.Latest,err}; defer os.RemoveAll(tmp)
	pkg := filepath.Join(tmp,assetName); sums := filepath.Join(tmp,"SHA256SUMS")
	if err=download(ctx,packageURL,pkg,1<<30);err!=nil{return false,info.Latest,err}
	if err=download(ctx,sumsURL,sums,4<<20);err!=nil{return false,info.Latest,err}
	want,err:=checksumFor(sums,assetName);if err!=nil{return false,info.Latest,err}
	got,err:=sha256File(pkg);if err!=nil{return false,info.Latest,err};if !strings.EqualFold(got,want){return false,info.Latest,errors.New("release package SHA-256 verification failed")}
	stage:=filepath.Join(tmp,"stage");if err=os.MkdirAll(stage,0o700);err!=nil{return false,info.Latest,err}
	if runtime.GOOS=="windows" {err=extractZip(pkg,stage)} else {err=extractTarGz(pkg,stage)};if err!=nil{return false,info.Latest,err}
	cfg,err:=config.Load(configPath);if err!=nil{return false,info.Latest,err}
	exe,err:=os.Executable();if err!=nil{return false,info.Latest,err}; binDir:=filepath.Dir(exe)
	backup:=filepath.Join(tmp,"backup"); _=os.MkdirAll(backup,0o700)
	if err=stopServer();err!=nil{return false,info.Latest,fmt.Errorf("stop server: %w",err)}
	rolledBack:=false
	defer func(){ if rolledBack { _=startServer() } }()
	if err=replaceWithBackup(filepath.Join(stage,"bin"),binDir,backup);err!=nil{_ = startServer();return false,info.Latest,err}
	webBackup:=filepath.Join(backup,"web")
	if err=copyDir(cfg.WebRoot,webBackup);err!=nil && !os.IsNotExist(err){rollback(binDir,backup,cfg.WebRoot,webBackup);_ = startServer();return false,info.Latest,err}
	if err=os.RemoveAll(cfg.WebRoot);err==nil{err=copyDir(filepath.Join(stage,"web"),cfg.WebRoot)}
	if err!=nil{rollback(binDir,backup,cfg.WebRoot,webBackup);_ = startServer();return false,info.Latest,err}
	if err=startServer();err!=nil{rollback(binDir,backup,cfg.WebRoot,webBackup);_ = startServer();return false,info.Latest,fmt.Errorf("restart failed, rollback attempted: %w",err)}
	if !waitHealth(){rollback(binDir,backup,cfg.WebRoot,webBackup);_ = startServer();return false,info.Latest,errors.New("updated server failed health check; rollback attempted")}
	return true,info.Latest,nil
}

func packageName(ver string) string { ext:="tar.gz"; if runtime.GOOS=="windows"{ext="zip"}; arch:=runtime.GOARCH; return fmt.Sprintf("sego-server_%s_%s_%s.%s",normalize(ver),runtime.GOOS,arch,ext) }
func download(ctx context.Context,url,path string,max int64)error{req,_:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil);req.Header.Set("User-Agent","SEGO-Server-Updater");resp,err:=(&http.Client{Timeout:5*time.Minute}).Do(req);if err!=nil{return err};defer resp.Body.Close();if resp.StatusCode!=200{return fmt.Errorf("download HTTP %d",resp.StatusCode)};f,err:=os.OpenFile(path,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0o600);if err!=nil{return err};defer f.Close();n,err:=io.Copy(f,io.LimitReader(resp.Body,max+1));if err!=nil{return err};if n>max{return errors.New("download exceeded size limit")};return nil}
func checksumFor(path,name string)(string,error){b,err:=os.ReadFile(path);if err!=nil{return "",err};for _,line:=range strings.Split(string(b),"\n"){f:=strings.Fields(line);if len(f)>=2&&strings.TrimPrefix(f[1],"*")==name{return f[0],nil}};return "",errors.New("release checksum entry missing")}
func sha256File(path string)(string,error){f,err:=os.Open(path);if err!=nil{return "",err};defer f.Close();h:=sha256.New();if _,err=io.Copy(h,f);err!=nil{return "",err};return hex.EncodeToString(h.Sum(nil)),nil}
func safeJoin(root,name string)(string,error){name=filepath.Clean(filepath.FromSlash(name));if filepath.IsAbs(name)||name==".."||strings.HasPrefix(name,".."+string(filepath.Separator)){return "",errors.New("unsafe archive path")};p:=filepath.Join(root,name);rel,err:=filepath.Rel(root,p);if err!=nil||rel==".."||strings.HasPrefix(rel,".."+string(filepath.Separator)){return "",errors.New("archive escape blocked")};return p,nil}
func extractZip(path,dst string)error{z,err:=zip.OpenReader(path);if err!=nil{return err};defer z.Close();for _,f:=range z.File{p,err:=safeJoin(dst,f.Name);if err!=nil{return err};if f.FileInfo().IsDir(){if err=os.MkdirAll(p,0o755);err!=nil{return err};continue};if err=os.MkdirAll(filepath.Dir(p),0o755);err!=nil{return err};r,err:=f.Open();if err!=nil{return err};w,err:=os.OpenFile(p,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0o755);if err==nil{_,err=io.Copy(w,r);w.Close()};r.Close();if err!=nil{return err}};return nil}
func extractTarGz(path,dst string)error{f,err:=os.Open(path);if err!=nil{return err};defer f.Close();gz,err:=gzip.NewReader(f);if err!=nil{return err};defer gz.Close();tr:=tar.NewReader(gz);for{h,err:=tr.Next();if err==io.EOF{break};if err!=nil{return err};p,err:=safeJoin(dst,h.Name);if err!=nil{return err};switch h.Typeflag{case tar.TypeDir:err=os.MkdirAll(p,0o755);case tar.TypeReg:if err=os.MkdirAll(filepath.Dir(p),0o755);err==nil{var w *os.File;w,err=os.OpenFile(p,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,os.FileMode(h.Mode)&0o777);if err==nil{_,err=io.Copy(w,tr);w.Close()}};default:continue};if err!=nil{return err}};return nil}
func replaceWithBackup(src,dst,backup string)error{ents,err:=os.ReadDir(src);if err!=nil{return err};for _,e:=range ents{if e.IsDir()||strings.Contains(strings.ToLower(e.Name()),"updater"){continue};sp:=filepath.Join(src,e.Name());dp:=filepath.Join(dst,e.Name());bp:=filepath.Join(backup,e.Name());if _,err:=os.Stat(dp);err==nil{_ = os.Rename(dp,bp)};if err=os.Rename(sp,dp);err!=nil{return err}};return nil}
func rollback(binDir,backup,web,webBackup string){ents,_:=os.ReadDir(backup);for _,e:=range ents{if e.IsDir(){continue};_ = os.Remove(filepath.Join(binDir,e.Name()));_ = os.Rename(filepath.Join(backup,e.Name()),filepath.Join(binDir,e.Name()))};_ = os.RemoveAll(web);_ = copyDir(webBackup,web)}
func copyDir(src,dst string)error{st,err:=os.Stat(src);if err!=nil{return err};if !st.IsDir(){return errors.New("source is not a directory")};return filepath.Walk(src,func(p string,info os.FileInfo,err error)error{if err!=nil{return err};rel,_:=filepath.Rel(src,p);out:=filepath.Join(dst,rel);if info.IsDir(){return os.MkdirAll(out,info.Mode().Perm())};in,err:=os.Open(p);if err!=nil{return err};defer in.Close();w,err:=os.OpenFile(out,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,info.Mode().Perm());if err!=nil{return err};_,err=io.Copy(w,in);cerr:=w.Close();if err!=nil{return err};return cerr})}
func stopServer()error{if runtime.GOOS=="windows"{return exec.Command("schtasks.exe","/End","/TN","SEGO Server").Run()};return exec.Command("systemctl","stop","sego.service").Run()}
func startServer()error{if runtime.GOOS=="windows"{return exec.Command("schtasks.exe","/Run","/TN","SEGO Server").Run()};return exec.Command("systemctl","start","sego.service").Run()}
func waitHealth()bool{c:=&http.Client{Timeout:2*time.Second};for i:=0;i<20;i++{r,err:=c.Get("http://127.0.0.1:8088/api/health");if err==nil{r.Body.Close();if r.StatusCode==200{return true}};time.Sleep(500*time.Millisecond)};return false}
