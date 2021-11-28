import { Box, Card } from "@mui/material";

export const ImageView = ({
  src,
  invalidated,
  invalidate,
}: {
  src?: string;
  invalidated: boolean;
  invalidate: () => void;
}) => {
  const content =
    src && !invalidated ? (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={src}
        alt="capture"
        style={{ width: "100%", height: "100%" }}
        onError={invalidate}
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
    );
  return (
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
          {content}
        </Box>
      </Box>
    </Card>
  );
};
