namespace GrasshopperMcp.Plugin.Grasshopper;

internal sealed class ComponentCatalog
{
    public object ListComponents()
    {
        return new
        {
            components = Array.Empty<object>()
        };
    }
}
