import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { TooltipProvider } from "@/components/ui/tooltip";
import { Header } from "@/components/layout/Header";
import { Footer } from "@/components/layout/Footer";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Scout9 - Automated Scouting Reports",
  description: "AI-powered scouting report generator for League of Legends and VALORANT esports teams. Built with GRID API data.",
  keywords: ["esports", "scouting", "League of Legends", "VALORANT", "analytics", "Cloud9"],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased min-h-screen flex flex-col bg-background`}
      >
        <TooltipProvider>
          <Header />
          <main className="flex-1">
            {children}
          </main>
          <Footer />
        </TooltipProvider>
      </body>
    </html>
  );
}
