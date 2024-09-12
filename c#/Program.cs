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
            var response = generative.CreateSession(7);
            string sessionId = response.Result.Id;

            // Sending a new message
            // You can implements file send inside this method
            var messageResponse = generative.SendMessage(sessionId, "Quero blusas e canecas");

            foreach (var message in messageResponse.Result.Response)
            {
                switch (message.Type)
                {
                    case "trigger":
                        Console.WriteLine("trigger:");
                        var dataJson = JsonConvert.SerializeObject(message.Data);
                        var sessionResponse = JsonConvert.DeserializeObject<TriggerResponse>(dataJson);
                        if (sessionResponse.FunctionName == "consulta_produto")
                        {
                            // Write your query to get product and send to user
                            Console.WriteLine("Consulta produto:");
                            if (sessionResponse.Args.ContainsKey("produto"))
                            {
                                Console.WriteLine(sessionResponse.Args["produto"]);
                            }

                            continue;
                        }
                        
                        if (sessionResponse.FunctionName == "adicionar_carrinho")
                        {
                            // Write your add product here
                            Console.WriteLine("Adicione ao carrinho:");
                            if (sessionResponse.Args.ContainsKey("produto"))
                            {
                                Console.WriteLine(sessionResponse.Args["produto"]);
                            }

                            continue;
                        }
                        
                        if (sessionResponse.FunctionName == "finalizar_compra")
                        {
                            // Write your finalize here
                            Console.WriteLine("Compra finalizada!");
                            continue;
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