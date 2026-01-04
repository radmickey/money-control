// Common API response helpers

/**
 * Extracts data from API response, handling both nested and flat response formats
 * API can return: { data: { data: {...} } } or { data: {...} }
 */
export const extractData = <T = any>(response: { data?: { data?: T } | T }): T | undefined => {
  const data = response.data as any;
  return data?.data !== undefined ? data.data : data;
};

/**
 * Extracts error message from API error response
 */
export const extractErrorMessage = (error: any, fallbackMessage: string): string => {
  return (
    error.response?.data?.error?.message ||
    error.response?.data?.error ||
    error.message ||
    fallbackMessage
  );
};

