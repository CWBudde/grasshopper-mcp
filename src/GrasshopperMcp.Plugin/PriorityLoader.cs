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
            new GraphMutationService()));

    public override GH_LoadingInstruction PriorityLoad()
    {
        Server.Start();
        RhinoApp.WriteLine($"Grasshopper MCP adapter listening on 127.0.0.1:{LocalServer.DefaultPort}.");
        return GH_LoadingInstruction.Proceed;
    }
}
