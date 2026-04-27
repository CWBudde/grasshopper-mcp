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
                    version = VersionInfo.Version,
                    activeDocument = _documents.HasActiveDocument,
                    grasshopperLoaded = true
                }),
            "document_info" => ProtocolResponse.Success(request.Id, _documents.GetDocumentInfo()),
            "list_components" => ProtocolResponse.Success(request.Id, _components.ListComponents()),
            "run_solver" => ProtocolResponse.Success(request.Id, _graph.RunSolver()),
            "add_component" => _graph.AddComponent(request.Id),
            "set_input" => _graph.SetInput(request.Id),
            "connect" => _graph.Connect(request.Id),
            "get_output" => _graph.GetOutput(request.Id),
            _ => ProtocolResponse.Failure(request.Id, "unknown_method", $"Unknown method '{request.Method}'.")
        };
    }
}
