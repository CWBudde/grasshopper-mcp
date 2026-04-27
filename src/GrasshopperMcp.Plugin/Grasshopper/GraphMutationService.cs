namespace GrasshopperMcp.Plugin.Grasshopper;

using GrasshopperMcp.Plugin.Server;

internal sealed class GraphMutationService
{
    public object RunSolver()
    {
        return new
        {
            completed = false,
            message = "Grasshopper document mutation is not implemented yet."
        };
    }

    public ProtocolResponse AddComponent(string requestId)
    {
        return NotImplemented(requestId, "add_component");
    }

    public ProtocolResponse SetInput(string requestId)
    {
        return NotImplemented(requestId, "set_input");
    }

    public ProtocolResponse Connect(string requestId)
    {
        return NotImplemented(requestId, "connect");
    }

    public ProtocolResponse GetOutput(string requestId)
    {
        return NotImplemented(requestId, "get_output");
    }

    private static ProtocolResponse NotImplemented(string requestId, string method)
    {
        return ProtocolResponse.Failure(
            requestId,
            "graph_mutation_not_implemented",
            $"The '{method}' command is defined but still needs live Grasshopper API wiring.");
    }
}
