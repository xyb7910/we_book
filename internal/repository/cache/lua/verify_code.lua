---
--- Created by ypb.
--- DateTime: 2024/7/21 21:14
---
local key = KEY[1]
-- 用户输入的 code
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
-- 转成一个数字
local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 用户一直输入错误
    return -1
elseif expectedCode == code then
    -- 用户输入正确
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户手抖输入错误
    redis.call("decr", cntKey)
    return -2
end