import 'package:flutter/material.dart';

import 'placeholder_screen.dart';

/// Root widget for the VEGA client.
///
/// Owns app-wide rendering concerns (theming, routing entry point) and the
/// initial screen. All async startup is performed by the host app before this
/// widget is mounted.
class VegaApp extends StatelessWidget {
  const VegaApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'VEGA',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
      ),
      home: const PlaceholderScreen(),
    );
  }
}
