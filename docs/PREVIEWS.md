# File Preview Security

SEGO previews files only after authorization and only for explicitly allowed types.

## Supported MVP previews

- images: JPEG, PNG, WebP, GIF, BMP, ICO
- PDF
- plain text / Markdown / JSON / CSV / logs
- selected source/configuration text formats
- audio: MP3, WAV, OGG, M4A, AAC, FLAC
- video: MP4, WebM, OGV, MOV, M4V

Actual audio/video playback depends on browser codec support.

## Preview flow

1. The authenticated user requests a preview token for a file in their namespace.
2. The server verifies the path, ownership scope, file type and preview size policy.
3. SEGO returns a temporary opaque token scoped to that file.
4. The browser loads `/api/preview/{token}` without putting the main bearer token in the media URL.
5. Video/audio is streamed with Range request support for seeking.

Media tokens have a longer temporary window than document/image tokens because browsers can make additional Range requests during normal playback.

## Active-content protection

HTML, HTM and SVG are not admitted as inline preview types in the MVP. Preview responses also receive restrictive content-security policy and `nosniff` protection.

## Emergency behavior

Emergency SOS deletes every active preview token and blocks preview access while lockdown remains active.
