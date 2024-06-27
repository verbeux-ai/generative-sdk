using System.Collections.Generic;
using Newtonsoft.Json;

namespace suri_demo
{
    public class SessionResponse
    {
        public string Id { get; set; }
        public int AssistantId { get; set; }
        public bool RestrictedByContext { get; set; }
        public string Description { get; set; }
        public bool IsActive { get; set; }
        // Adicione outras propriedades conforme necessário
    }
    
    public class SessionUpdateResponse
    {
        public string Id { get; set; }
        public bool IsAnythingElse { get; set; }
        public List<SessionUpdateResponseBodyResponse> Response { get; set; }
    }

    public class SessionUpdateResponseBodyResponse
    {
        public string Type { get; set; }
        public object Data { get; set; }
    }
    
    public class TriggerResponse
    {
        [JsonProperty("function_name")]
        public string FunctionName { get; set; }

        [JsonProperty("args")]
        public Dictionary<string, string> Args { get; set; }
    }
}