import funtch from 'funtch';

/**
 * Return configuration variables for app
 * @return {Object} Configuration object
 */
export default async function getConfig() {
  if (process.env.NODE_ENV === 'production') {
    return await funtch.get('/env');
  }

  return {
    HERODOTE_API: process.env.REACT_APP_HERODOTE_API,
    ALGOLIA_APP: process.env.REACT_APP_ALGOLIA_APP,
    ALGOLIA_KEY: process.env.REACT_APP_ALGOLIA_KEY,
    ALGOLIA_INDEX: process.env.REACT_APP_ALGOLIA_INDEX,
  };
}
