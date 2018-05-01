import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

let response200StatusCounter = new Counter("response_200_status");
let response500StatusCounter = new Counter("response_500_status");
let responseOtherStatusCounter = new Counter("response_other_status");
let numberOfServices = 4;
let init = new Map();

export default function() {
  let res = http.get("http://load-balancer");

  let code = res.status;
  switch (code) {
    case 200:
      response200StatusCounter.add(1);
      break;
    case 500:
      response500StatusCounter.add(1);
      console.log("[ERROR] 500 " + res.body);
      break;
    default:
      responseOtherStatusCounter.add(1);
      console.log("[ERROR] " + res.status + " " + res.body);
  }

  check(res, {
    "status was 200": (r) => r.status === 200
  });

  sleep(1);
};

export function setup() {

  while (true) {
    let res = http.get("http://load-balancer");
    if ((res.status === 200) && !init.has(res.body)) {
      init.set(res.body, true);
      console.log("[INFO] Warm up found: " + res.body)
    }

    if (init.size === numberOfServices) {
      console.log("[INFO] Warm up completed");
      break;
    }
  }

}
