/// <reference types="react" />
/// <reference types="react-dom" />
/// <reference types="node" />

declare namespace NodeJS {
  interface ProcessEnv {
    NEXT_PUBLIC_API_URL: string;
    NEXT_PUBLIC_OAUTH_REDIRECT_URL: string;
    NEXT_PUBLIC_APP_NAME: string;
    NEXT_PUBLIC_APP_DESCRIPTION: string;
  }
}
