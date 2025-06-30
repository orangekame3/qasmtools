import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'QASM Tools Playground',
  description: 'Online playground for OpenQASM 3.0 formatting and validation',
  keywords: ['QASM', 'OpenQASM', 'quantum', 'formatter', 'playground'],
  authors: [{ name: 'orangekame3' }],
  openGraph: {
    title: 'QASM Tools Playground',
    description: 'Online playground for OpenQASM 3.0 formatting and validation',
    type: 'website',
    url: 'https://orangekame3.github.io/qasmtools/',
    images: [
      {
        url: '/og-image.png',
        width: 1200,
        height: 630,
        alt: 'QASM Tools Playground',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'QASM Tools Playground',
    description: 'Online playground for OpenQASM 3.0 formatting and validation',
    images: ['/og-image.png'],
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" data-theme="qasm">
      <body className={inter.className} suppressHydrationWarning={true}>
        <div className="min-h-screen bg-base-100">
          {children}
        </div>
      </body>
    </html>
  )
}