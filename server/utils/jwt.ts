import jwt from "jsonwebtoken";
import { useRuntimeConfig } from "#imports";

export interface JWTPayload {
  userId: number;
  email: string;
  roles: string[];
}

export const createJWT = (payload: JWTPayload): string => {
  const config = useRuntimeConfig();
  return jwt.sign(payload, config.sessionSecret, {
    expiresIn: "7d",
  });
};

export const verifyJWT = (token: string): JWTPayload | null => {
  try {
    const config = useRuntimeConfig();
    const decoded = jwt.verify(token, config.sessionSecret) as JWTPayload;
    return decoded;
  } catch {
    return null;
  }
};
