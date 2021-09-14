import React, { useEffect, useState } from "react";
import { Box, Button, Card, Slider } from "@material-ui/core";

interface HistoryResponse {
  files: [HistoryFile];
}

interface HistoryFile {
  path: string;
  timestamp: number;
  invalidated: boolean | undefined;
}

const CameraHistoryViewer: React.FC = () => {
  const [selected, setSelected] = useState(0);
  const [imageList, setImageList] = useState([] as HistoryFile[]);

  useEffect(() => {
    fetch("/api/history")
      .then((res) => res.json())
      .then(
        (res: HistoryResponse) => {
          setImageList(res.files);
          setSelected(res.files.length - 1);
        },
        (e) => {
          console.error(e);
          setImageList([]);
        }
      );
  }, []);

  const min = 0;
  const max = Math.max(imageList.length - 1, 0);

  const selectPrev = () => {
    setSelected((selected) => Math.max(min, selected - 1));
  };

  const selectNext = () => {
    setSelected((selected) => Math.min(max, selected + 1));
  };

  return (
    <>
      <Box sx={{ mb: 2 }}>
        {imageList[selected]
          ? new Date(imageList[selected].timestamp * 1000).toLocaleString()
          : "Loading..."}
      </Box>
      <Card>
        <Box
          style={{
            // aspectRatio: "16/9",
            paddingTop: "56.25%",
            position: "relative",
            width: "100%",
          }}
        >
          <Box
            style={{
              position: "absolute",
              top: "0",
              width: "100%",
              height: "100%",
            }}
          >
            {imageList[selected] && !imageList[selected].invalidated ? (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                src={imageList[selected]?.path}
                alt="capture"
                style={{ width: "100%", height: "100%" }}
                onError={() => {
                  const newList = [...imageList];
                  newList[selected] = {
                    ...newList[selected],
                    invalidated: true,
                  };
                  setImageList(newList);
                }}
              />
            ) : (
              <Box
                width="100%"
                height="100%"
                display="flex"
                alignItems="center"
                justifyContent="center"
                style={{
                  backgroundColor: "#bbb",
                  color: "#fff",
                  height: "100%",
                  fontSize: "5em",
                }}
                sx={{ m: 0 }}
              >
                ?
              </Box>
            )}
          </Box>
        </Box>
      </Card>
      <Box display="flex" justifyContent="center">
        <Button
          size="large"
          variant="outlined"
          sx={{ m: 1 }}
          onClick={selectPrev}
        >
          &lt;
        </Button>
        <Button
          size="large"
          variant="outlined"
          sx={{ m: 1 }}
          onClick={selectNext}
        >
          &gt;
        </Button>
      </Box>
      <Slider
        sx={{ my: 2 }}
        disabled={imageList.length == 0}
        step={1}
        min={min}
        max={max}
        value={selected}
        onChange={(e, v) => setSelected(v as number)}
      />
    </>
  );
};
export { CameraHistoryViewer };
