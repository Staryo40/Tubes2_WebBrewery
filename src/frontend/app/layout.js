export const metadata = {
  title: "Graph Visualizer",
  description: "Visualize BFS & DFS Traversals",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}