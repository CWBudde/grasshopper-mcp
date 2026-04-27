using GrasshopperMcp.Plugin.Grasshopper;

namespace GrasshopperMcp.Plugin.Server;

internal sealed class CommandRouter
{
    private readonly DocumentService _documents;
    private readonly ComponentCatalog _components;
    private readonly GraphMutationService _graph;

    public CommandRouter(DocumentService documents, ComponentCatalog components, GraphMutationService graph)
    {
        _documents = documents;
        _components = components;
        _graph = graph;
    }

    public ProtocolResponse Route(ProtocolRequest request)
    {
        return request.Method switch
        {
            "health" => ProtocolResponse.Success(request.Id, new
                {
                    version = "0.1.0",
                    activeDocument = _documents.HasActiveDocument,
                    grasshopperLoaded = true
                }),
            "document_info" => ProtocolResponse.Success(request.Id, _documents.GetDocumentInfo()),
            "list_components" => ProtocolResponse.Success(request.Id, _components.ListComponents()),
            "run_solver" => ProtocolResponse.Success(request.Id, _graph.RunSolver()),
            _ => ProtocolResponse.Failure(request.Id, "unknown_method", $"Unknown method '{request.Method}'.")
        };
    }
}
