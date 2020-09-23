# Copyright 2020 Yoshi Yamaguchi
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import urllib
from locust import HttpUser, task, between
from hypothesis import given
from hypothesis import strategies
from random import choice

servers = ["34.84.62.35", "35.243.116.90", "34.84.17.51"]
port = "8080"
source = strategies.text(min_size=5, max_size=100)


class Client(HttpUser):
    host = "http://" + choice(servers) + ":" + port
    wait_time = between(0.5, 2)

    @task(1)
    def request(self):
        query = source.example()
        params = urllib.parse.urlencode({"q": query})
        self.client.get("/?" + params)