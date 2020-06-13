import funtch from "funtch";

/**
 * Return configuration variables for app
 * @return {Object} Configuration object
 */
export default async function getConfig() {
  if (process.env.NODE_ENV === "production") {
    return await funtch.get("/env");
  }

  return {
    ALGOLIA_APPLICATION_ID: process.env.REACT_APP_ALGOLIA_APPLICATION_ID,
    ALGOLIA_API_KEY: process.env.REACT_APP_ALGOLIA_API_KEY,
    ALGOLIA_INDEX: process.env.REACT_APP_ALGOLIA_INDEX,
  };
}
