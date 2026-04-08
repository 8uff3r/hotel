import bcrypt from "bcrypt";
import { verifyJWT, type JWTPayload } from "../utils/jwt";

const publicPaths = ["/api/auth/login", "/api/auth/logout", "/api/auth/me"];

const rolePermissions: Record<string, string[]> = {
  "/api/users": ["admin"],
  "/api/rooms": ["admin", "manager", "receptionist"],
  "/api/guests": ["admin", "manager", "receptionist"],
  "/api/reservations": ["admin", "manager", "receptionist"],
  "/api/accounts": ["admin", "manager"],
  "/api/expenses": ["admin", "manager"],
  "/api/income": ["admin", "manager"],
  "/api/parking": ["admin", "manager", "receptionist", "staff"],
};

export default defineEventHandler((event) => {
  const path = event.path;
  const method = event.method;

  if (!path.startsWith("/api")) {
    return;
  }

  const isPublic = publicPaths.some((p) => path.startsWith(p));
  if (isPublic) {
    return;
  }

  const authToken = getCookie(event, "auth_token");
  if (!authToken) {
    throw createError({
      statusCode: 401,
      statusMessage: "Unauthorized: No token provided",
    });
  }

  const payload = verifyJWT(authToken);
  if (!payload) {
    throw createError({
      statusCode: 401,
      statusMessage: "Unauthorized: Invalid token",
    });
  }

  event.context.user = payload as JWTPayload;

  const requiredRoles = getRequiredRoles(path);
  if (requiredRoles.length > 0) {
    const userRoles = (payload as JWTPayload).roles || [];
    const hasPermission = requiredRoles.some((role) => userRoles.includes(role));

    if (!hasPermission) {
      throw createError({
        statusCode: 403,
        statusMessage: "Forbidden: Insufficient permissions",
      });
    }
  }
});

function getRequiredRoles(path: string): string[] {
  for (const [routePath, roles] of Object.entries(rolePermissions)) {
    if (path.startsWith(routePath)) {
      return roles;
    }
  }
  return [];
}
