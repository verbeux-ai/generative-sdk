using System;
using Newtonsoft.Json;

namespace suri_demo
{
    internal class Program
    {
        public static void Main(string[] args)
        {
            // Instantiates a generative
            // Put your api key here (https://admin.verbeux.com.br/api-keys)
            Generative generative = new Generative("...-...-...-...-....");
            
            // Creating a new session
            // Get your assistant_id in admin 
            var response = generative.CreateSession(9999); // 
            string sessionId = response.Result.Id;

            // Sending a new message
            // You can implements file send inside this method
            var messageResponse = generative.SendMessage(sessionId, "Quero falar com uma atendente");

            foreach (var message in messageResponse.Result.Response)
            {
                switch (message.Type)
                {
                    case "trigger":
                        Console.WriteLine("trigger:");
                        var dataJson = JsonConvert.SerializeObject(message.Data);
                        var sessionResponse = JsonConvert.DeserializeObject<TriggerResponse>(dataJson);
                        if (sessionResponse.FunctionName == "direcionar_atendente")
                        {
                            // Triggered 
                            // Redirects to human service or other
                        }
                        
                        Console.WriteLine(dataJson);
                        break;
                    case "text":
                        Console.WriteLine("text:");
                        var textJson = JsonConvert.SerializeObject(message.Data);
                        Console.WriteLine(textJson);
                        break;
                    case "reference":
                        Console.WriteLine("references:");
                        var referenceJson = JsonConvert.SerializeObject(message.Data);
                        Console.WriteLine(referenceJson);
                        break;
                }
            }
        }
    }
}