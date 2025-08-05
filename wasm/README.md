# LBaaS Packet Catcher - WebAssembly Edition

This directory contains the WebAssembly version of the LBaaS Packet Catcher game, built with Go and Ebiten.

## Quick Start

### Option 1: Using Go HTTP Server (Recommended)

```bash
# From the project root
make wasm-serve-go
```

Then open http://localhost:8000/ in your browser.

### Option 2: Using wasmserve (Alternative)

```bash
# From the project root
make wasm-serve
```

Then open http://localhost:8080/ in your browser.

### Option 3: Manual Setup

```bash
# Build WebAssembly version
make wasm-setup

# Serve with Python HTTP server
make wasm-serve-simple
```

Then open http://localhost:8000/ in your browser.

## Files

- `lbbaspack.wasm` - The compiled WebAssembly binary
- `wasm_exec.js` - Go's WebAssembly runtime (copied from Go installation)
- `index.html` - The HTML page that loads and runs the game

## Browser Compatibility

The WebAssembly version works in all modern browsers that support WebAssembly:

- Chrome 57+
- Firefox 52+
- Safari 11+
- Edge 16+

## Controls

- **Mouse**: Move the load balancer to catch packets
- **Space**: Start the game
- **Escape**: Return to menu or exit
- **Ctrl+X**: Exit the game

## Troubleshooting

### "wasm_exec.js not found" Error

If you get an error about `wasm_exec.js` not being found:

1. Make sure you have Go installed
2. Run `make wasm-copy-js` to copy the file from your Go installation
3. If the file still can't be found, download it manually from:
   - Go 1.24+: https://raw.githubusercontent.com/golang/go/go1.24.x/lib/wasm/wasm_exec.js
   - Go 1.23-: https://raw.githubusercontent.com/golang/go/go1.23.x/misc/wasm/wasm_exec.js

### "Failed to load WebAssembly" Error

1. Make sure you're serving the files from an HTTP server (not opening the HTML file directly)
2. Check that all files (`lbbaspack.wasm`, `wasm_exec.js`, `index.html`) are in the same directory
3. Check the browser console for more detailed error messages

### MIME Type Errors

If you see errors like "MIME type mismatch" or "was blocked due to MIME type":

1. **Use the Go HTTP server**: `make wasm-serve-go` (recommended)
2. **Or configure your web server** to serve the correct MIME types:
   - `.wasm` files: `application/wasm`
   - `.js` files: `application/javascript`

The Go HTTP server automatically sets the correct MIME types.

### Audio Context Error

If you see "The AudioContext was not allowed to start" error:

1. This is normal - the game will work without audio
2. If you want audio, click anywhere on the page first to enable audio context

### Fullscreen Error

If you see "Fullscreen request denied" error:

1. **Click the "🎮 Fullscreen" button** below the game canvas
2. Fullscreen must be triggered by a user gesture (click) in modern browsers
3. The game works perfectly in windowed mode as well

### Canvas Display Issues

If the game appears to draw outside the canvas:

1. **Click the "▶️ Start Game" button** to focus the canvas
2. The canvas is properly sized at 800x600 pixels
3. The game container has a green border to show the game area
4. **Only one canvas is used** - Ebiten's canvas is automatically styled and positioned

### Two Canvas Issue (Fixed)

Previously, there were two canvases - one created by Ebiten and one in our HTML. This has been fixed by:

1. Using a placeholder div instead of a canvas in HTML
2. Letting Ebiten create its own canvas
3. Moving Ebiten's canvas to our game container
4. Applying our styling to Ebiten's canvas

### Performance Issues

1. Make sure you're using a modern browser
2. Close other tabs/applications to free up memory
3. The game runs at 60 FPS by default

## Development

### Rebuilding the WebAssembly Binary

```bash
make wasm
```

### Testing Changes

```bash
# Clean and rebuild
make clean
make wasm-setup
make wasm-serve
```

### Debugging

1. Open browser developer tools (F12)
2. Check the Console tab for error messages
3. Use the Network tab to verify files are loading correctly

## Deployment

To deploy the WebAssembly version:

1. Run `make wasm-setup` to create all necessary files
2. Upload the contents of the `wasm/` directory to your web server
3. Ensure your web server serves `.wasm` files with the correct MIME type:
   - Content-Type: `application/wasm`

### Example Nginx Configuration

```nginx
location ~* \.wasm$ {
    add_header Content-Type application/wasm;
}

location ~* \.js$ {
    add_header Content-Type application/javascript;
}
```

### Example Apache Configuration

Add to `.htaccess`:
```apache
AddType application/wasm .wasm
AddType application/javascript .js
```

## Architecture

The WebAssembly version uses the same ECS (Entity-Component-System) architecture as the native version:

- **Entities**: Game objects (load balancer, backends, packets, etc.)
- **Components**: Data containers (transform, sprite, collider, etc.)
- **Systems**: Logic processors (movement, collision, rendering, etc.)

All game logic runs in the browser using WebAssembly, providing near-native performance.

## Performance Notes

- Initial load time may be 1-3 seconds depending on connection speed
- The game runs at 60 FPS on most modern devices
- Memory usage is typically 50-100MB
- The WebAssembly binary is approximately 5-10MB

## Browser Security

The WebAssembly version runs in a sandboxed environment and cannot:
- Access the file system
- Make network requests (except for the initial load)
- Access system resources outside the browser

This makes it safe to run in any modern web browser. 