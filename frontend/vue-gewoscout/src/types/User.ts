/**
 * User interface represents the structure of a logged in user object.
 */
export interface User {
    identityProvider: string;
    userId: string;
    userDetails: string;
    userRoles: string[];
}