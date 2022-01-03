import React from "react";
import { Button } from "@mui/material";

const OpenLiveButton: React.FC = () => {
  const onClickOpenLive = () => {
    fetch("/api/streamUri")
      .then((res) => res.json())
      .then((res) => {
        const streamUri = res.uri;
        document.location.href = `vlc-x-callback://x-callback-url/stream?url=${streamUri}`;
      });
  };
  return <Button onClick={onClickOpenLive}>Open Live</Button>;
};
export { OpenLiveButton };
