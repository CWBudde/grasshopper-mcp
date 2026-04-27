using Grasshopper.Kernel;

namespace GrasshopperMcp.Plugin;

public sealed class GrasshopperMcpAssemblyInfo : GH_AssemblyInfo
{
    public override string Name => "Grasshopper MCP";
    public override string Description => "Local Grasshopper adapter for the Go grasshopper-mcp server.";
    public override string AuthorName => "grasshopper-mcp contributors";
    public override string AuthorContact => "";
    public override string Version => "0.1.0";
    public override GH_LibraryLicense License => GH_LibraryLicense.opensource;
}
