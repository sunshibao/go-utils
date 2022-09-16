package redis

import (
	"github.com/garyburd/redigo/redis"
)

const SCRIPT_VOICE2 = `
local voiceGiftNumKey = tostring(KEYS[1])
local voiceTemKey = tostring(KEYS[2])
local giftCom = tostring(ARGV[1])
local gnum = tonumber(ARGV[2])
local guest = tonumber(ARGV[3])
local tmpInfomarshal = tostring(ARGV[4])
local keyValidity = tonumber(ARGV[5])

if giftCom ~= 'A' and giftCom ~= 'B' and giftCom ~= 'C'
then 
  return 3
end

redis.call('hincrby',voiceGiftNumKey,giftCom,gnum)

local giftNumArray = redis.call('HMGET',voiceGiftNumKey,'A','B','C')

local minGiftNum = 99999
for i = 1,3 do
  local tmpNum = 0
  if giftNumArray[i] ~= false
  then
    tmpNum = tonumber(giftNumArray[i])
  end

  if tmpNum <= 0
  then
    return 2
  else
    if tmpNum< minGiftNum
    then
      minGiftNum = tmpNum
    end
  end
end

local newArray={}
for i = 1,3 do
  local tmpNum2 = tonumber(giftNumArray[i])
  newArray[i]=tmpNum2-minGiftNum
end
redis.call('HMSET',voiceGiftNumKey,'A',newArray[1],'B',newArray[2],'C',newArray[3])
if keyValidity > 0
then
	local voice= redis.call('SET',voiceTemKey,tmpInfomarshal,'EX',keyValidity,'NX')
	if voice ~= false
	then
		return 1
	end
end
return 0
`

var voiceScript2 = redis.NewScript(2, SCRIPT_VOICE2)

func main() {

	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return
	}
	defer c.Close()

	reply, err := redis.Int(voiceScript2.Do(c, "Voice:RoomID:{686}:SeatId:{7c871778e2c011ea98b600163e012514}:TemInfo", "voiceGiftNum:{86}:{86}", "A", 1, 86))

	if reply == 0 {
		return
	}
}
