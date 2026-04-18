import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const backendUrl = process.env.BACKEND_URL;

  if (backendUrl && request.nextUrl.pathname.startsWith("/api/")) {
    const target = new URL(
      request.nextUrl.pathname + request.nextUrl.search,
      backendUrl
    );
    return NextResponse.rewrite(target);
  }
}

export const config = {
  matcher: "/api/:path*",
};
