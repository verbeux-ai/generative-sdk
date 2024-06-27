using System;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace suri_demo
{
    public class Generative
    {
        private string _baseUrl;
        private string _apiKey;

        public Generative(string apiKey, string baseUrl = "https://generative-api.verbeux.com.br")
        {
            _baseUrl = baseUrl;
            _apiKey = apiKey;
        }

        // CreateSession creates a session based on assistantId
        // Check https://generative-api.verbeux.com.br/docs/index.html#/sess%C3%A3o/post_session_
        public async Task<SessionResponse> CreateSession(int assistantId)
        {
            using (HttpClient client = new HttpClient())
            {
                client.DefaultRequestHeaders.Add("api-key", _apiKey);

                var requestBody = new
                {
                    assistant_id = assistantId
                };

                var json = JsonConvert.SerializeObject(requestBody);
                var content = new StringContent(json, Encoding.UTF8, "application/json");

                HttpResponseMessage response = await client.PostAsync($"{_baseUrl}/session/", content);

                if (response.IsSuccessStatusCode)
                {
                    var responseBody = await response.Content.ReadAsStringAsync();
                    var sessionResponse = JsonConvert.DeserializeObject<SessionResponse>(responseBody);
                    return sessionResponse;
                }

                throw new Exception($"Failed to create session: {response.StatusCode} - {response.ReasonPhrase}");
            }
        }

        // SendMessage sends a message to an existing session and returns the response from the assistant
        public async Task<SessionUpdateResponse> SendMessage(string sessionId, string message)
        {
            using (HttpClient client = new HttpClient())
            {
                client.DefaultRequestHeaders.Add("api-key", _apiKey);

                var requestBody = new
                {
                    message
                };

                var json = JsonConvert.SerializeObject(requestBody);
                var content = new StringContent(json, Encoding.UTF8, "application/json");

                HttpResponseMessage response = await client.PutAsync($"{_baseUrl}/session/{sessionId}", content);

                if (response.IsSuccessStatusCode)
                {
                    var responseBody = await response.Content.ReadAsStringAsync();
                    var sessionUpdateResponse = JsonConvert.DeserializeObject<SessionUpdateResponse>(responseBody);
                    return sessionUpdateResponse;
                }

                throw new Exception(
                    $"Failed to send message: {response.StatusCode} - {response.ReasonPhrase}: {response.Content.ReadAsStringAsync()}");
            }
        }
    }
}