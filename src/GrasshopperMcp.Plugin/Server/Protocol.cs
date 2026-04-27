using System.Text.Json;

namespace GrasshopperMcp.Plugin.Server;

internal sealed record ProtocolRequest(string Id, string Method, JsonElement? Params);

internal sealed record ProtocolResponse(string Id, bool Ok, object? Result, ProtocolError? Error)
{
    public static ProtocolResponse Success(string id, object? result) => new(id, true, result, null);

    public static ProtocolResponse Failure(string id, string code, string message) =>
        new(id, false, null, new ProtocolError(code, message));
}

internal sealed record ProtocolError(string Code, string Message);
