import "../styles/globals.css";
import type { AppProps } from "next/app";
import Head from "next/dist/shared/lib/head";
import React from "react";
import { ThemeProvider, CssBaseline } from "@mui/material";
import { blue, pink } from "@mui/material/colors";
import { createTheme } from "@mui/material/styles";

function MyApp({ Component, pageProps }: AppProps) {
  const [darkMode, setDarkMode] = React.useState(false);

  const theme = createTheme({
    palette: {
      primary: blue,
      secondary: pink,
      mode: darkMode ? "dark" : "light",
    },
  });

  React.useEffect(() => {
    if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
      setDarkMode(true);
    } else {
      setDarkMode(false);
    }
  }, []);
  return (
    <>
      <Head>
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no"
        />
      </Head>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <Component {...pageProps} />
      </ThemeProvider>
    </>
  );
}
export default MyApp;
