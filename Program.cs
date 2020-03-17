using Topshelf;

namespace ButtonBoxVjoyFeeder
{
    public class Program
    {
        public static void Main(string[] args)
        {
            HostFactory.Run(config =>
            {
                config.Service<SerialPortVJoyFeederService>(service =>
                {
                    service.ConstructUsing(() => new SerialPortVJoyFeederService());
                    service.WhenStarted(x => x.Start());
                    service.WhenStopped(x => x.Stop());
                });

                config.RunAsLocalSystem();
                config.StartAutomatically();
                config.SetDisplayName(nameof(ButtonBoxVjoyFeeder));
                config.SetServiceName(nameof(ButtonBoxVjoyFeeder));
            });
        }
    }
}