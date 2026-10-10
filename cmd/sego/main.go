package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nadeemmhdm/sego-server/internal/auth"
	"github.com/nadeemmhdm/sego-server/internal/config"
	segodb "github.com/nadeemmhdm/sego-server/internal/db"
	"github.com/nadeemmhdm/sego-server/internal/securitylock"
	"github.com/nadeemmhdm/sego-server/internal/storage"
	"github.com/nadeemmhdm/sego-server/internal/updater"
	"github.com/nadeemmhdm/sego-server/internal/version"
	"golang.org/x/term"
)

const dashboardURL = "http://127.0.0.1:8088"

func main() {
	cfgPath := flag.String("config", defaultConfigPath(), "config path")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 { quickHelp(); return }

	joined := strings.Join(args," ")
	switch joined {
	case "help", "--help", "-h": usage(); return
	case "version": fmt.Println("SEGO Server", version.Version); return
	case "open": openDashboard(); return
	case "start": fatalIf(controlService("start")); fmt.Println("[OK] SEGO Server started"); return
	case "stop": fatalIf(controlService("stop")); fmt.Println("[OK] SEGO Server stopped"); return
	case "restart": fatalIf(controlService("restart")); fmt.Println("[OK] SEGO Server restarted"); return
	case "status": showSimpleStatus(); return
	case "tunnel setup": fatalIf(tunnelSetup()); return
	case "tunnel status": tunnelStatus(); return
	case "tunnel remove": fatalIf(tunnelRemove()); return
	case "update check": updateCheck(); return
	case "update auto": fatalIf(runUpdater(*cfgPath)); return
	}

	cfg, err := config.Ensure(*cfgPath); fatalIf(err)
	if len(args)==2 && args[0]=="mode" {
		switch args[1] {
		case "show": fmt.Printf("Mode: %s\nUpload concurrency: %d\n",config.Mode(cfg),cfg.UploadConcurrency); return
		case config.ModeLowResource, config.ModeStandard:
			fatalIf(config.ApplyMode(&cfg,args[1])); fatalIf(config.Save(*cfgPath,cfg)); fmt.Printf("[OK] mode set to %s\n",args[1]); return
		}
	}
	if len(args)==3 && args[0]=="mode" && args[1]=="set" { fatalIf(config.ApplyMode(&cfg,args[2])); fatalIf(config.Save(*cfgPath,cfg)); fmt.Printf("[OK] mode set to %s\n",args[2]); return }

	d,err:=segodb.Open(cfg.DatabasePath);fatalIf(err);defer d.Close()
	switch joined {
	case "setup": ensureAdmin(d); fmt.Println("\nSetup complete.\nOpen SEGO with: sego open")
	case "doctor": doctor(cfg)
	case "storage list": for _,r:=range cfg.StorageRoots{fmt.Printf("%s\t%s\t%s\n",r.ID,r.Name,r.Path)}
	case "admin create": createAdmin(d)
	case "admin ensure": ensureAdmin(d)
	case "emergency status": showEmergencyStatus(d)
	case "emergency unlock": unlockEmergency(d)
	default: usage()
	}
}

func defaultConfigPath() string { if v:=strings.TrimSpace(os.Getenv("SEGO_CONFIG"));v!=""{return v};if runtime.GOOS=="windows"{if pd:=os.Getenv("ProgramData");pd!=""{return filepath.Join(pd,"SEGO","sego.json")}};if _,err:=os.Stat("/etc/sego/sego.json");err==nil{return "/etc/sego/sego.json"};return "sego.json" }
func showSimpleStatus(){running:=healthOK();fmt.Println("SEGO Server");if running{fmt.Println("  Server    : running\n  Dashboard : "+dashboardURL)}else{fmt.Println("  Server    : not reachable")};if s:=serviceState();s!=""{fmt.Println("  Service   : "+s)};if !running{fmt.Println("\nTry: sego start")}}
func openDashboard(){if !healthOK(){fmt.Println("SEGO is not running. Starting it now...");fatalIf(controlService("start"));for i:=0;i<20&&!healthOK();i++{time.Sleep(500*time.Millisecond)}};if !healthOK(){fatalIf(fmt.Errorf("server did not become ready; run `sego status` and `sego doctor`"))};if err:=launchBrowser(dashboardURL);err!=nil{fmt.Println("Open this address in your browser:\n"+dashboardURL);return};fmt.Println("[OK] Opening "+dashboardURL)}
func healthOK()bool{c:=&http.Client{Timeout:1200*time.Millisecond};r,err:=c.Get(dashboardURL+"/api/health");if err!=nil{return false};defer r.Body.Close();return r.StatusCode==200}
func serviceState()string{if runtime.GOOS=="windows"{if exec.Command("schtasks.exe","/Query","/TN","SEGO Server").Run()==nil{return "installed (starts at boot)"};return "not installed"};if runtime.GOOS=="linux"{out,err:=exec.Command("systemctl","is-active","sego.service").Output();if err==nil{return strings.TrimSpace(string(out))+" (boot enabled)"};return "inactive"};return ""}
func controlService(action string)error{if runtime.GOOS=="windows"{var err error;switch action{case"start":err=runCommand("schtasks.exe","/Run","/TN","SEGO Server");case"stop":err=runCommand("schtasks.exe","/End","/TN","SEGO Server");case"restart":_ = runCommand("schtasks.exe","/End","/TN","SEGO Server");time.Sleep(500*time.Millisecond);err=runCommand("schtasks.exe","/Run","/TN","SEGO Server")};if err==nil{return nil};exe,e:=os.Executable();if e!=nil{return err};ps:=fmt.Sprintf("$p=Start-Process -FilePath '%s' -ArgumentList '%s' -Verb RunAs -Wait -PassThru; exit $p.ExitCode",strings.ReplaceAll(exe,"'","''"),action);return runCommand("powershell.exe","-NoProfile","-ExecutionPolicy","Bypass","-Command",ps)};if runtime.GOOS=="linux"{a:=[]string{action,"sego.service"};if runCommand("systemctl",a...)==nil{return nil};return runCommand("sudo",append([]string{"systemctl"},a...)...)};return fmt.Errorf("service control unsupported on %s",runtime.GOOS)}
func runCommand(name string,args ...string)error{c:=exec.Command(name,args...);c.Stdout=os.Stdout;c.Stderr=os.Stderr;c.Stdin=os.Stdin;return c.Run()}
func launchBrowser(url string)error{switch runtime.GOOS{case"windows":return exec.Command("rundll32","url.dll,FileProtocolHandler",url).Start();case"linux":return exec.Command("xdg-open",url).Start();case"darwin":return exec.Command("open",url).Start()};return fmt.Errorf("browser launch unsupported")}

func updateCheck(){ctx,cancel:=context.WithTimeout(context.Background(),20*time.Second);defer cancel();info,err:=updater.Latest(ctx,version.Version);fatalIf(err);if info.Latest==""{fmt.Println("No published SEGO release is available yet.");return};if info.Available{fmt.Printf("Update available: %s -> %s\nRun: sego update\n",info.Current,info.Latest)}else{fmt.Printf("SEGO %s is up to date.\n",info.Current)}}
func runUpdater(cfg string)error{exe,err:=os.Executable();if err!=nil{return err};name:="sego-updater";if runtime.GOOS=="windows"{name+=".exe"};u:=filepath.Join(filepath.Dir(exe),name);if _,err=os.Stat(u);err!=nil{return fmt.Errorf("updater is not installed; rerun the SEGO installer")};if runtime.GOOS=="linux"&&os.Geteuid()!=0{return runCommand("sudo",u,"--config",cfg)};if runtime.GOOS=="windows"{ps:=fmt.Sprintf("$p=Start-Process -FilePath '%s' -ArgumentList '--config','%s' -Verb RunAs -Wait -PassThru; exit $p.ExitCode",strings.ReplaceAll(u,"'","''"),strings.ReplaceAll(cfg,"'","''"));return runCommand("powershell.exe","-NoProfile","-ExecutionPolicy","Bypass","-Command",ps)};return runCommand(u,"--config",cfg)}

func tunnelSetup()error{
	fmt.Println("Cloudflare Tunnel setup\nBefore continuing, create a remotely-managed tunnel in Cloudflare and set its public hostname service to http://localhost:8088.")
	token:=strings.TrimSpace(readSecret("Paste the tunnel token (hidden): "))
	if len(token)<20{return fmt.Errorf("tunnel token looks invalid")}
	if runtime.GOOS=="linux"{
		if _,err:=exec.LookPath("cloudflared");err!=nil{
			fmt.Println("Installing cloudflared from Cloudflare's signed APT repository...")
			cmd:=`set -e; sudo mkdir -p --mode=0755 /usr/share/keyrings; curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null; echo "deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared any main" | sudo tee /etc/apt/sources.list.d/cloudflared.list >/dev/null; sudo apt-get update; sudo apt-get install -y cloudflared`
			if err:=runCommand("bash","-lc",cmd);err!=nil{return err}
		}
		_ = runCommand("sudo","cloudflared","service","uninstall")
		if err:=runCommand("sudo","cloudflared","service","install",token);err!=nil{return err}
		recoverCmd:=`set -e; sudo mkdir -p /etc/systemd/system/cloudflared.service.d; printf '%s\n' '[Unit]' 'After=network-online.target sego.service' 'Wants=network-online.target' '' '[Service]' 'Restart=always' 'RestartSec=5s' | sudo tee /etc/systemd/system/cloudflared.service.d/sego-restart.conf >/dev/null; sudo systemctl daemon-reload; sudo systemctl enable --now cloudflared; sudo systemctl restart cloudflared`
		if err:=runCommand("bash","-lc",recoverCmd);err!=nil{return err}
		if err:=exec.Command("systemctl","is-active","--quiet","cloudflared").Run();err!=nil{return fmt.Errorf("cloudflared service was installed but is not running")}
		fmt.Println("[OK] Tunnel connected. It starts at boot and automatically reconnects after network or process failures.")
		return nil
	}
	if runtime.GOOS=="windows"{
		ps:=`$ErrorActionPreference='Stop'; if(-not(Get-Command cloudflared.exe -ErrorAction SilentlyContinue)){winget install --id Cloudflare.cloudflared -e --silent --accept-package-agreements --accept-source-agreements | Out-Host}; $cf=(Get-Command cloudflared.exe -ErrorAction SilentlyContinue).Source; if(-not $cf){$cf=(Get-ChildItem 'C:\Program Files','C:\Program Files (x86)' -Filter cloudflared.exe -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1 -ExpandProperty FullName)}; if(-not $cf){throw 'cloudflared.exe not found after install'}; & $cf service uninstall 2>$null | Out-Null; & $cf service install $env:SEGO_TUNNEL_TOKEN; sc.exe config cloudflared start= auto | Out-Null; sc.exe failure cloudflared reset= 0 actions= restart/5000/restart/15000/restart/60000 | Out-Null; sc.exe failureflag cloudflared 1 | Out-Null; $svc=Get-Service cloudflared -ErrorAction Stop; if($svc.Status -ne 'Running'){Start-Service cloudflared}; Start-Sleep -Seconds 1; if((Get-Service cloudflared).Status -ne 'Running'){throw 'cloudflared service did not start'}`
		c:=exec.Command("powershell.exe","-NoProfile","-ExecutionPolicy","Bypass","-Command",ps)
		c.Env=append(os.Environ(),"SEGO_TUNNEL_TOKEN="+token);c.Stdout=os.Stdout;c.Stderr=os.Stderr;c.Stdin=os.Stdin
		if err:=c.Run();err!=nil{return err}
		fmt.Println("[OK] Tunnel connected. The Windows service starts at boot and automatically restarts after failures.")
		return nil
	}
	return fmt.Errorf("tunnel setup unsupported on %s",runtime.GOOS)
}
func tunnelStatus(){if runtime.GOOS=="windows"{_ = runCommand("sc.exe","query","cloudflared");return};if runtime.GOOS=="linux"{_ = runCommand("systemctl","--no-pager","status","cloudflared");return};fmt.Println("Tunnel status unsupported on this OS")}
func tunnelRemove()error{if runtime.GOOS=="windows"{return runCommand("powershell.exe","-NoProfile","-Command","$c=(Get-Command cloudflared.exe -ErrorAction Stop).Source; & $c service uninstall")};if runtime.GOOS=="linux"{_ = runCommand("sudo","rm","-f","/etc/systemd/system/cloudflared.service.d/sego-restart.conf");_ = runCommand("sudo","systemctl","daemon-reload");return runCommand("sudo","cloudflared","service","uninstall")};return fmt.Errorf("tunnel removal unsupported")}

func ensureAdmin(d *sql.DB){var n int;fatalIf(d.QueryRow(`SELECT COUNT(*) FROM users WHERE role IN ('super_admin','admin')`).Scan(&n));if n>0{fmt.Printf("[OK] administrator already configured (%d admin account(s))\n",n);return};fmt.Println("\nCreate your SEGO administrator account\nThis account controls users, security and server settings.");createAdmin(d)}
func createAdmin(d *sql.DB){r:=bufio.NewReader(os.Stdin);fmt.Print("Admin username: ");u,_:=r.ReadString('\n');u=strings.TrimSpace(u);if len(u)<3{fatalIf(fmt.Errorf("username must be at least 3 characters"))};p1:=readSecret("Password (12+ characters): ");p2:=readSecret("Confirm password: ");if p1!=p2{fatalIf(fmt.Errorf("passwords do not match"))};h,err:=auth.HashPassword(p1);fatalIf(err);_,err=d.Exec(`INSERT INTO users(username,password_hash,role,created_at) VALUES(?,?,?,?)`,u,h,"super_admin",time.Now().UTC().Format(time.RFC3339));fatalIf(err);fmt.Println("[OK] Super Admin created")}
func showEmergencyStatus(d *sql.DB){s,err:=securitylock.StateNow(context.Background(),d);fatalIf(err);if !s.Active{fmt.Println("Emergency lockdown: inactive");return};fmt.Printf("Emergency lockdown: ACTIVE\nActivated at: %s\nReason: %s\n",s.ActivatedAt,s.Reason)}
func unlockEmergency(d *sql.DB){state,err:=securitylock.StateNow(context.Background(),d);fatalIf(err);if !state.Active{fmt.Println("[OK] Emergency lockdown is already inactive");return};r:=bufio.NewReader(os.Stdin);fmt.Print("Admin username: ");u,_:=r.ReadString('\n');u=strings.TrimSpace(u);p:=readSecret("Admin password: ");var id int64;var h,role string;err=d.QueryRow(`SELECT id,password_hash,role FROM users WHERE username=?`,u).Scan(&id,&h,&role);if err!=nil||(role!="super_admin"&&role!="admin")||!auth.VerifyPassword(h,p){fatalIf(fmt.Errorf("admin verification failed"))};fatalIf(securitylock.DeactivateLocal(context.Background(),d,id,u,role));fmt.Println("[OK] Emergency lockdown disabled.\nAll users must sign in again. Revoked share links remain revoked for safety.")}
func doctor(cfg config.Config){roots:=make([]storage.Root,0,len(cfg.StorageRoots));for _,r:=range cfg.StorageRoots{roots=append(roots,storage.Root{ID:r.ID,Name:r.Name,Path:r.Path})};_,err:=storage.New(roots);fatalIf(err);fmt.Println("[OK] configuration\n[OK] database path configured\n[OK] storage roots accessible");fmt.Printf("[INFO] mode: %s\n",config.Mode(cfg));if healthOK(){fmt.Println("[OK] local API reachable")}else{fmt.Println("[WARN] local API is not reachable at "+dashboardURL)}}
func readSecret(prompt string)string{fmt.Print(prompt);if term.IsTerminal(int(os.Stdin.Fd())){b,err:=term.ReadPassword(int(os.Stdin.Fd()));fmt.Println();if err==nil{return string(b)}};s,_:=bufio.NewReader(os.Stdin).ReadString('\n');return strings.TrimSpace(s)}
func fatalIf(err error){if err!=nil{fmt.Fprintln(os.Stderr,"SEGO:",err);os.Exit(1)}}
func quickHelp(){fmt.Println(`SEGO Server

Easy commands:
  sego               Open the dashboard
  sego status        Check the server
  sego start         Start SEGO
  sego stop          Stop SEGO
  sego remote        Set up remote access once
  sego update        Install the latest GitHub Release update
  sego update check  Check for a release update
  sego help          Show commands`)}
func usage(){fmt.Println(`SEGO Server CLI

Everyday:
  sego open | status | start | stop | restart

Setup / diagnostics:
  sego setup
  sego doctor
  sego storage list

Remote tunnel:
  sego tunnel setup
  sego tunnel status
  sego tunnel remove

Updates:
  sego version
  sego update check
  sego update auto

Emergency:
  sego emergency status
  sego emergency unlock

Performance:
  sego mode show
  sego mode low-resource
  sego mode standard

Advanced custom config:
  sego --config /path/to/sego.json <command>`)}
