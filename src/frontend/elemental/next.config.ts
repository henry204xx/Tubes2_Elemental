import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  eslint: {
    ignoreDuringBuilds: true, // ✅ Ignores ESLint errors during production build (e.g., on Vercel)
  },
  // Add other config options if needed
};

export default nextConfig;
