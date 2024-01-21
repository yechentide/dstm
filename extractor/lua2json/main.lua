if #arg < 1 then
    print("Usage: main.lua <directory> <filename>")
    os.exit(1)
end

local directory = arg[1]
local filename = arg[2]
local filename_without_extension = filename:gsub("%..+$", "")

package.path = package.path..';'..directory..'/?.lua;'

require("my_utils")
local luaObject = require(filename_without_extension)

local outputPath = directory.."/"..filename_without_extension..".json"
local jsonStr = ItemToJson(luaObject, 0)
WriteToFile(outputPath, jsonStr)
