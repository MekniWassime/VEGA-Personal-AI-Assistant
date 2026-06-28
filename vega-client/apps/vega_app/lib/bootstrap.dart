import 'package:flutter/widgets.dart';

/// Boots the app.
///
/// The app layer owns all startup work — especially async initialization
/// (config, storage, service clients, …) — so that the `core` package can stay
/// focused purely on rendering. [builder] supplies the root widget once
/// initialization has finished.
Future<void> bootstrap(Widget Function() builder) async {
  WidgetsFlutterBinding.ensureInitialized();

  await _initialize();

  runApp(builder());
}

/// Async startup work. Extend this as services are added.
Future<void> _initialize() async {
  // TODO: real async initialization (env, storage, service clients, …).
}
