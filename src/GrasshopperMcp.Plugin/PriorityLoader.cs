using Grasshopper.Kernel;
using GrasshopperMcp.Plugin.Grasshopper;
using GrasshopperMcp.Plugin.Server;
using Rhino;

namespace GrasshopperMcp.Plugin;

public sealed class PriorityLoader : GH_AssemblyPriority
{
    private static readonly LocalServer Server = new(
        new CommandRouter(
            new DocumentService(),
            new ComponentCatalog(),
            new GraphMutationService()),
        ReadPort());

    public override GH_LoadingInstruction PriorityLoad()
    {
        try
        {
            Server.Start();
            RhinoApp.WriteLine($"Grasshopper MCP adapter {VersionInfo.Version} listening on 127.0.0.1:{Server.Port}.");
            return GH_LoadingInstruction.Proceed;
        }
        catch (Exception ex)
        {
            RhinoApp.WriteLine($"Grasshopper MCP adapter failed to start on 127.0.0.1:{Server.Port}: {ex.Message}");
            return GH_LoadingInstruction.Abort;
        }
    }

    private static int ReadPort()
    {
        var raw = Environment.GetEnvironmentVariable("GRASSHOPPER_MCP_PORT");
        return int.TryParse(raw, out var port) && port > 0 && port <= 65535
            ? port
            : LocalServer.DefaultPort;
    }
}
