$source = @"
using System;
public class McBedrockTool {
    public static string GetKeyByCoords(int x, int z, int y = 0, int Dimension = 0, byte Tag = 0x2f) {
        // C# doesn't seem to handle negative integer division well, so have to divide a double then floor it
        int ChunkX = (int)Math.Floor(x / 16.0);
        int ChunkZ = (int)Math.Floor(z / 16.0);
        byte SubChunkY = (byte) (y / 16);
        Console.WriteLine(String.Format("Chunk indices:\n X: {0} ({1})\nZ: {2} ({3})\nY: {4} ({5})", ChunkX, x, ChunkZ, z, (y / 16), y));
        string MyKey = HexKey(ChunkX) +
            HexKey(ChunkZ) +
            (Dimension == 0 ? "" : HexKey(Dimension)) +
            HexKey(Tag) +
            (Tag == 0x2f ? HexKey(SubChunkY) : "");
        return MyKey;
    }
    public static string HexKey(int i) {
        byte[] ByteArray = BitConverter.GetBytes(i);
        // Force byte order to little endian no matter the local platform
        if (! BitConverter.IsLittleEndian)
            Array.Reverse(ByteArray);
        return String.Concat(Array.ConvertAll(ByteArray, x => x.ToString("X2")));
    }
    public static string HexKey(byte i) {
        return i.ToString("X2");
    }
}
"@

Add-Type -TypeDefinition $source 
$env:MCPETOOL_WORLD="C:\Users\danha\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds\XfILYVNgAQA="

# 43 68 -8 : warp wood planks
$key = [McBedrockTool]::GetKeyByCoords($args[0], $args[1], $args[2])
Write-Host coords: $args[0] $args[1] $args[2]
Write-Host mcpetool db get $key
#mcpetool db get $key


#PS C:\Users\danha\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds> Write-Host $env:MCPETOOL_WORLD; .\getdata.ps1 "43" "-8" "68" | sls minecraft:
#C:\Users\danha\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds\ex79YOL+AAA=
#
#          "value": "minecraft:air"
#          "value": "minecraft:sand"
#          "value": "minecraft:grass"
#          "value": "minecraft:dirt"
#          "value": "minecraft:leaves"
#          "value": "minecraft:double_plant"
#          "value": "minecraft:double_plant"
#          "value": "minecraft:tallgrass"
#          "value": "minecraft:yellow_flower"
#          "value": "minecraft:red_flower"
#          "value": "minecraft:warped_planks"
#

#"Overworld subchunk coordinate"
#[McBedrockTool]::GetKeyByCoords(413,54,90)
## 19000000030000002F05
#
#"Negative example"
#[McBedrockTool]::GetKeyByCoords(-413,54,90)
## E6FFFFFF030000002F05
#
#"Nether dimension"
#[McBedrockTool]::GetKeyByCoords(-413,54,90, -1)
## E6FFFFFF03000000FFFFFFFF2F05
#
#"End dimension"
#[McBedrockTool]::GetKeyByCoords(413,-54,90, 1)
## 19000000FCFFFFFF010000002F05
#
#"Overworld block entity data key"
## Y is irrelevant here and Dimension redundant, but I can't use named parameters from PowerShell to C#
#[McBedrockTool]::GetKeyByCoords(413,54,90, 0, 50)
## 190000000300000032