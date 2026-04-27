using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Text.Json;

namespace GrasshopperMcp.Plugin.Server;

internal sealed class LocalServer : IDisposable
{
    public const int DefaultPort = 47820;

    private static readonly JsonSerializerOptions JsonOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.CamelCase,
        PropertyNameCaseInsensitive = true
    };

    private readonly CommandRouter _router;
    private TcpListener? _listener;
    private CancellationTokenSource? _stop;
    private Task? _acceptLoop;

    public int Port { get; }

    public LocalServer(CommandRouter router, int port = DefaultPort)
    {
        _router = router;
        Port = port;
    }

    public void Start()
    {
        if (_listener is not null)
        {
            return;
        }

        _stop = new CancellationTokenSource();
        _listener = new TcpListener(IPAddress.Loopback, Port);
        _listener.Start();
        _acceptLoop = Task.Run(() => AcceptLoopAsync(_stop.Token));
    }

    public void Stop()
    {
        _stop?.Cancel();
        _listener?.Stop();
        _listener = null;
    }

    public void Dispose()
    {
        Stop();
        _stop?.Dispose();
    }

    private async Task AcceptLoopAsync(CancellationToken cancellationToken)
    {
        while (!cancellationToken.IsCancellationRequested && _listener is not null)
        {
            try
            {
                var client = await _listener.AcceptTcpClientAsync(cancellationToken).ConfigureAwait(false);
                _ = Task.Run(() => HandleClientAsync(client, cancellationToken), cancellationToken);
            }
            catch (OperationCanceledException)
            {
                return;
            }
            catch (ObjectDisposedException)
            {
                return;
            }
        }
    }

    private async Task HandleClientAsync(TcpClient client, CancellationToken cancellationToken)
    {
        using (client)
        {
            await using var stream = client.GetStream();
            using var reader = new StreamReader(stream, Encoding.UTF8, leaveOpen: true);
            await using var writer = new StreamWriter(stream, new UTF8Encoding(false), leaveOpen: true)
            {
                AutoFlush = true,
                NewLine = "\n"
            };

            var line = await reader.ReadLineAsync(cancellationToken)
                .ConfigureAwait(false);
            if (string.IsNullOrWhiteSpace(line))
            {
                return;
            }

            ProtocolResponse response;
            try
            {
                var request = JsonSerializer.Deserialize<ProtocolRequest>(line, JsonOptions);
                response = request is null
                    ? ProtocolResponse.Failure("", "invalid_request", "Request body was empty.")
                    : _router.Route(request);
            }
            catch (JsonException ex)
            {
                response = ProtocolResponse.Failure("", "invalid_json", ex.Message);
            }
            catch (Exception ex)
            {
                response = ProtocolResponse.Failure("", "internal_error", ex.Message);
            }

            var responseJson = JsonSerializer.Serialize(response, JsonOptions);
            await writer.WriteLineAsync(responseJson).ConfigureAwait(false);
        }
    }
}
