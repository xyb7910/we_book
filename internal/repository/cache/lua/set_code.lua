---
--- Created by ypb.
--- DateTime: 2024/7/21 21:15
---

-- 验证码我们保存在 redis 的 key 中
-- phone_code:login:
local key = KEYS[1]
-- 验证次数，一个验证码，最多重复三次，这个纪录还可以验证几次
local cntKey = key..":cnt"
-- 你的验证码为 123456
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在，但是没有过期时间
   return -2
elseif ttl == -2 or ttl < 540 then
   redis.call("set", key, val)
   redis.call("expire", key, 600)
   redis.call("set", cntKey, 3)
   redis.call("expire", cntKey, 600)
   -- 符合我们的预期
   return 0
else
   -- 发送太频繁
   return -1
end