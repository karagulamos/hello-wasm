const https = require("https");

exports.handler = async function (event, context) {
  const credentials = process.env.USERNAME + ":" + process.env.PASSWORD;
  const token = Buffer.from(credentials).toString("base64");

  const response = await doHttpRequest(
    {
      hostname: "broker-api.sandbox.alpaca.markets",
      path: event.path,
      headers: { Authorization: `Basic ${token}` },
    },
    Uint8Array.from([])
  );

  return {
    statusCode: response.statusCode,
    body: JSON.stringify({
      name: response.data.name,
      symbol: response.data.symbol,
    }),
  };
};

function doHttpRequest(options, data) {
  return new Promise((resolve, reject) => {
    const request = https.request(options, (response) => {
      response.setEncoding("utf8");
      let responseBody = "";

      response.on("data", (chunk) => {
        responseBody += chunk;
      });

      response.on("end", () => {
        resolve({
          statusCode: response.statusCode,
          data: JSON.parse(responseBody),
        });
      });
    });

    request.on("error", (err) => {
      reject(err);
    });

    request.write(data);
    request.end();
  });
}
