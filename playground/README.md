# QASM Tools Playground

An interactive web playground for OpenQASM 3.0 formatting and validation, powered by WebAssembly.

## Features

- **Real-time Formatting**: Format OpenQASM 3.0 code instantly in your browser
- **Split-pane Interface**: Input on the left, formatted output on the right
- **Monaco Editor**: Full-featured code editor with syntax highlighting
- **WebAssembly Powered**: Fast formatting using the same engine as the CLI tool
- **Mobile Responsive**: Works on desktop and mobile devices
- **No Server Required**: Runs entirely in your browser

## Usage

1. **Input Code**: Type or paste your OpenQASM 3.0 code in the left panel
2. **Format**: Click the "Format" button to format your code
3. **Copy Output**: Use the "Copy" button to copy the formatted code to clipboard
4. **Load Example**: Click "Load Example" to try a sample QASM program
5. **Clear**: Use "Clear" to start with an empty editor

## Development

### Prerequisites

- Node.js 18+ 
- Go 1.22+

### Setup

```bash
# Install dependencies
cd playground
npm install

# Build WebAssembly module (from project root)
cd ..
task build:wasm

# Copy WASM files to public directory
cp -r bin/wasm playground/public/

# Start development server
cd playground
npm run dev
```

### Build for Production

```bash
# Build static site
npm run build

# The output will be in the 'out' directory
```

## Deployment

The playground is automatically deployed to GitHub Pages when changes are pushed to the main branch. See `.github/workflows/deploy-playground.yml` for the deployment configuration.

## Architecture

- **Frontend**: Next.js with TypeScript and daisyUI
- **Backend**: Go WebAssembly module
- **Editor**: Monaco Editor for syntax highlighting and editing
- **Hosting**: GitHub Pages with static site generation

## Browser Compatibility

- Chrome 57+
- Firefox 52+
- Safari 11+
- Edge 16+

WebAssembly support is required for the formatting functionality.

## Contributing

1. Make changes to the playground or WASM module
2. Test locally with `npm run dev`
3. Submit a pull request
4. Changes will be automatically deployed on merge

## License

Same as the main qasmtools project - see LICENSE in the project root.