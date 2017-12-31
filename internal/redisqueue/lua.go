package redisqueue

var priorityQueueScript = `
-- redis-priority-queue
-- Author: Gabriel Bordeaux (gabfl)
-- Github: https://github.com/gabfl/redis-priority-queue
-- Version: 1.0.4
-- (can only be used in 3.2+)

-- Get mandatory vars
local action = ARGV[1];
local queueName = ARGV[2];

-- returns true if empty or null
-- http://stackoverflow.com/a/19667498/50501
local function isempty(s)
  return s == nil or s == '' or type(s) == 'userdata'
end

-- Making sure required fields are not nil
assert(not isempty(action), 'ERR1: Action is missing')
assert(not isempty(queueName), 'ERR2: Queue name is missing')

if action == 'push'
then
     -- Define vars
    local item = ARGV[3];
    local priority = ARGV[4] or 100;

    -- Making sure required fields are not nil
    assert(not isempty(item), 'ERR5: Item is missing')

    -- Add item to queue
    return redis.call('ZADD', queueName, 'NX', priority, item)
elseif action == 'pop'
then
    -- Retrieve items
    local popped = redis.call('ZREVRANGEBYSCORE', queueName, '+inf', '-inf', 'LIMIT', 0, '1')
    if popped then
        for _,item in ipairs(popped) do
            -- Remove item
            redis.call('ZREM', queueName, item)
            return item;
        end
    end
return nil;
elseif action == 'count'
then
    -- Define vars
    local fromMin = ARGV[3] or '-inf';
    local toMax = ARGV[4] or '+inf';

    -- return queue count
    local count = redis.call('ZCOUNT', queueName, fromMin, toMax)

    return count;
else
    error('ERR3: Invalid action.')
end
`
