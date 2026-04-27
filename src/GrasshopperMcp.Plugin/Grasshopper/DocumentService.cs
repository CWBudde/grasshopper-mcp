namespace GrasshopperMcp.Plugin.Grasshopper;

internal sealed class DocumentService
{
    public bool HasActiveDocument => false;

    public object GetDocumentInfo()
    {
        return new
        {
            documentName = "",
            objectCount = 0,
            hasActiveDocument = HasActiveDocument
        };
    }
}
