const express = require("express");
const { networkInterfaces } = require("os");
const onvif = require("onvif");
const port = 3333;

const getLocalIpAddress = () => {
  const nets = networkInterfaces();
  const results = {};

  for (const name of Object.keys(nets)) {
    for (const net of nets[name]) {
      if (net.family === "IPv4" && !net.internal) {
        if (!results[name]) {
          results[name] = [];
        }
        results[name].push(net.address);
      }
    }
  }

  if (results["en0"] != null) {
    return results["en0"][0];
  }
  if (results["wlan0"] != null) {
    return results["wlan0"][0];
  }
  const keys = Object.keys(results);
  if (results[keys[0]] != null) {
    return results[keys[0]][0];
  }
  return null;
};

const getCameraRtspUri = async () => {
  return new Promise((resolve, reject) => {
    console.log("Discovering camera ...");
    onvif.Discovery.probe((err, cams) => {
      // function will be called only after timeout (5 sec by default)
      if (err) {
        reject(err);
        return;
      }
      console.log(`${cams.length} camera(s) found.`);
      if (cams.length == 0) {
        reject(new Error("Camera not found."));
        return;
      }
      console.log("Using first camera.");
      const cam = cams[0];
      console.log(`Hostname: ${cam.hostname}`);

      console.log("Getting stream uri ...");
      cam.getStreamUri({ protocol: "RTSP" }, (err, stream) => {
        if (err != null) {
          reject(err);
          return;
        }
        console.log("Receive stream uri.");
        console.log(stream);
        if (stream.uri == null) {
          reject(new Error("Uri is null."));
          return;
        }
        resolve(stream.uri);
      });
    });
  });
};

const app = express();
app.use(express.static("public"));

app.get("/", (req, res) => {
  req.url = "/index.html";
  app.handle(req, res);
});

app.get("/streamUri", (req, res) => {
  (async () => {
    try {
      const uri = await getCameraRtspUri();
      res.send({ uri, result: "OK" });
    } catch (err) {
      res.status(500).send({
        result: "NG",
        message: err.message,
        stack: err.stack,
      });
    }
  })();
});

const serverIpAddress = getLocalIpAddress();

app.listen(port, () => {
  console.log(`Server started. http://${serverIpAddress}:${port}`);
});
