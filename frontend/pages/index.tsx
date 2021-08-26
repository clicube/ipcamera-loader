import {Container} from "@material-ui/core";
import type {NextPage} from "next";
import Head from "next/dist/shared/lib/head";
import React from "react";
import {CameraHistoryViewer} from "../components/CameraHistoryViewer";
import {OpenLiveButton} from "../components/OpenLiveButton";

const Home: NextPage = () => {
  return (
    <>
      <Head>
        <meta name="apple-mobile-web-app-title" content="カメラを見る" />
        <link rel="apple-touch-icon-precomposed" href="icon.png" />
        <title>IP Camera util</title>
      </Head>
      <Container maxWidth="md" sx={{p: 4}}>
        <CameraHistoryViewer />
        <OpenLiveButton />
      </Container>
    </>
  );
};

export default Home;
