var AWS = require("aws-sdk");
var { has } = require("lodash");

var cloudfront = new AWS.CloudFront();

exports.handler = async (event) => {
  if (!has(event, "queryStringParameters.cfid")) {
    return {
      statusCode: 200,
    };
  }

  // Parameters
  const cloudfrontDistributionId = event.queryStringParameters.cfid;

  var params = {
    DistributionId: cloudfrontDistributionId,
    InvalidationBatch: {
      CallerReference: Date.now().toString(),
      Paths: {
        Quantity: 1,
        Items: ["/*"],
      },
    },
  };

  try {
    const invalidation = await cloudfront.createInvalidation(params).promise();
    return {
      statusCode: 200,
      body: JSON.stringify(invalidation),
    };
  } catch (e) {
    return {
      statusCode: 404,
      body: JSON.stringify(e),
    };
  }
};
