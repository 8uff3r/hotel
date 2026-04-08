export default defineEventHandler(async (event) => {
  try {
    // Clear the JWT cookie
    deleteCookie(event, "auth_token", {
      path: "/",
    });

    return {
      message: "Logout successful",
    };
  } catch {
    throw createError({
      statusCode: 500,
      message: "An error occurred during logout",
    });
  }
});
