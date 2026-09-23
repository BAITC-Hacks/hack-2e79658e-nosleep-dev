import type { NextConfig } from "next";
import path from "node:path";

const nextConfig: NextConfig = {
  turbopack: { root: path.resolve(__dirname, "../..") },
  output: "standalone",
  outputFileTracingRoot: path.resolve(__dirname, "../.."),
  async rewrites() {
    return [{
      source: "/api/v1/:path*",
      destination: "http://127.0.0.1:8000/api/v1/:path*",
    }];
  },
};

export default nextConfig;
