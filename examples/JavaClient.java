import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.WebSocket;
import java.nio.ByteBuffer;
import java.util.concurrent.CompletionStage;
import org.json.JSONObject;

public class JavaClient {
    public static void main(String[] args) throws Exception {
        HttpClient client = HttpClient.newHttpClient();
        WebSocket.Listener listener = new WebSocket.Listener() {
            @Override
            public void onOpen(WebSocket webSocket) {
                System.out.println("WebSocket connection opened");
                webSocket.request(1);
            }

            @Override
            public CompletionStage<?> onText(WebSocket webSocket, CharSequence data, boolean last) {
                System.out.println("Received message: " + data);
                webSocket.request(1);
                return null;
            }

            @Override
            public void onError(WebSocket webSocket, Throwable error) {
                System.out.println("Error occurred: " + error.getMessage());
            }
        };

        WebSocket webSocket = client.newWebSocketBuilder()
                .buildAsync(URI.create("ws://localhost:7502/ws"), listener)
                .join();

        while (true) {
            JSONObject message = new JSONObject();
            message.put("from", "java-client");
            message.put("to", "server");
            message.put("subject", "Hello");
            message.put("content", "How are you?");
            message.put("type", "msg");

            webSocket.sendText(message.toString(), true);
            Thread.sleep(5000);
        }
    }
}