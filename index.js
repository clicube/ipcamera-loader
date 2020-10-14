const http = require("http");
const { networkInterfaces } = require("os");
const onvif = require("onvif");

function main() {
  process.on("uncaughtException", function (err) {
    console.log(err);
  });
  const serverIpAddress = getLocalIpAddress();
  startServer(serverIpAddress);
}

function getLocalIpAddress() {
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
  const keys = Object.keys(results);
  if (results[keys[0]] != null) {
    return results[keys[0]][0];
  }
  return null;
}

async function getCameraRtspUri() {
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
}

function startServer(serverIpAddress) {
  http
    .createServer(async (req, res) => {
      try {
        const rtspUri = await getCameraRtspUri();
        const uri = `vlc://${rtspUri}`
        res.writeHead(302, {
          Location: rtspUri,
        });
        res.end();
        console.log(`Sent 302 response. ${uri}`);
      } catch (err) {
        console.log(err);
        res.writeHead(500);
        res.write(err.toString());
        res.write("\n\n");
        res.write(err.message);
        res.write("\n");
        res.write(err.stack);
        res.end();
      }
    })
    .listen(3000, () => {
      console.log(`Server started. http://${serverIpAddress}:3000`);
    });
}

main();
