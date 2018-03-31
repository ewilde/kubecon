import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

export let options = {
  vus: 10,
  duration: "300s"
};

var response200StatusCounter = new Counter("response_200_status");
var response500StatusCounter = new Counter("response_500_status");
var responseOtherStatusCounter = new Counter("response_other_status");

export default function() {
  let res = http.get("http://load-balancer");

  let code = res.status;
  switch (code) {
    case 200:
      response200StatusCounter.add(1);
      break;
    case 500:
      response500StatusCounter.add(1);
      break;
    default:
      responseOtherStatusCounter.add(1);
  }

  check(res, {
    "status was 200": (r) => r.status === 200
  });

  sleep(1);
};
