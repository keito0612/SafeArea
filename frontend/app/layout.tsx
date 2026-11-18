import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "SafeArea - 災害リスク可視化サービス",
  description: "引っ越し前に災害リスクを確認。住所を入力するだけで、洪水・地震・土砂災害などのリスクが一目でわかります。",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
